package app

import com.typesafe.scalalogging.Logger
import org.joda.time.DateTime
import org.joda.time.format.DateTimeFormat
import scopt.OParser

import java.nio.file.Paths

object Main extends App {
  val logger = Logger("diary-generator")

  val parser = {
    val builder = OParser.builder[Config]
    import builder._
    OParser.sequence(
      opt[String]("base-directory-path").required().action((baseDirectoryPath, config) => config.copy(baseDirectoryPath = baseDirectoryPath)),
      cmd("init")
        .action((_, config) => config.copy(mode = "init"))
        .children(
          opt[String]("template-path").required().action((templatePath, config) => config.copy(templatePathOpt = Some(templatePath)))
        ),
      cmd("archive")
        .action((_, config) => config.copy(mode = "archive"))
        .children(
          opt[String]("starts-with").required().action((startsWith, config) => config.copy(startsWithOpt = Some(startsWith)))
        )
    )
  }

  def init(baseDirectoryPath: String, templateFilePath: String, dt: DateTime): Unit = {
    val dateTimeFormatter = DateTimeFormat.forPattern("yyyy-MM-dd")
    val baseDir = Paths.get(baseDirectoryPath)

    val templateFile = Paths.get(templateFilePath).toFile

    val targetDirPath = baseDir.resolve(s"${dateTimeFormatter.print(dt)}").toAbsolutePath
    val targetDir = targetDirPath.toFile
    val targetFilePath = targetDirPath.resolve(templateFile.getName)
    val targetFile = targetFilePath.toFile

    if (templateFile.exists()) {
      if (targetDir.exists()) {
        logger.info(s"already exist : ${targetDirPath.toString}")
      } else {
        logger.info(s"make directory : ${targetDirPath.toString}")
        targetDir.mkdirs()
      }

      if (targetFile.exists()) {
        logger.info(s"already exist : ${targetFilePath.toString}")
      } else {
        logger.info(s"create file : ${targetFilePath.toString}")
        // todo テンプレートファイルに書式を設定できるようにする
        targetFile.createNewFile()
      }
    } else {
      logger.error("template file not found")
    }
  }

  OParser
    .parse(parser, args, Config())
    .fold {
      // do nothing
    } { config =>
      (config.mode, config.templatePathOpt, config.startsWithOpt) match {
        case ("init", Some(templateFilePath), _) =>
          init(config.baseDirectoryPath, templateFilePath, DateTime.now())
        case ("archive", _, Some(_)) => ??? // todo impl
        case _ => ??? // todo impl
      }
    }
}
