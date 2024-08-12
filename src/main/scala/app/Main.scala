package app

import org.joda.time.DateTime
import scopt.OParser

object Main extends App {

  private val parser = {
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

  OParser
    .parse(parser, args, Config())
    .fold {
      // do nothing
    } { config =>
      (config.mode, config.templatePathOpt, config.startsWithOpt) match {
        case ("init", Some(templateFilePath), _) =>
          DiaryInitializer.init(config.baseDirectoryPath, templateFilePath, DateTime.now)
        case ("archive", _, Some(startsWith)) =>
          DiaryArchiver.archive(config.baseDirectoryPath, startsWith)
        case _ => ??? // todo impl
      }
    }
}
